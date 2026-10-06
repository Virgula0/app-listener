# Graph Report - app-listener  (2026-10-06)

## Corpus Check
- 319 files · ~1,495,294 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4218 nodes · 16856 edges · 172 communities (139 shown, 33 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1236 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `71e70e1c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- testKeyPair
- monitor.bpf.c
- User
- guardUnitTest
- mkdirs
- Hash
- CandidateDir
- trustManager
- runUninstall
- highlight_test.go
- buildOneGuard
- .startGuardStd
- guard_trust.bpf.c
- GenerateConf
- fileEditModel
- New
- usecase/daemon_test.go
- Resource
- KernelDev
- Candidate
- guard.bpf.c
- Run
- NetEventType
- SanitizeText
- unlockUnderGuard
- SelectOrphans
- serve.go
- engine
- codeedit.go
- fileedit_model_test.go
- DaemonEvent
- daemonconfig.go
- update.go
- GuardEvent
- guarded_ancestor_within_limit
- sanitizeTerminalText
- .GrantTemporaryAccess
- IntegrationSuite
- dataRowRenderer
- main_test.go
- os.FileMode
- Monitor
- .syncSetsLocked
- IntegrationSuite
- filevault.go
- netModel
- IntegrationSuite
- runInstall
- github.com/charmbracelet/bubbletea.Cmd
- absPath
- os.File
- controlServer
- runForward
- EventType
- TrustGuard
- GuardInodeKey
- inode_dev
- fakeGuardRepo
- string
- pinstate.go
- LibraryClosure
- time.Duration
- github.com/spf13/cobra.Command
- catalogRefresher
- copyRegularAt
- bufio.Reader
- auditGroupCoverage
- btrfsLayoutFrom
- filevault_test.go
- serveWebSocket
- IntegrationSuite
- planUpdaters
- process_vm_readv.c
- runOfflineEdit
- systemd.go
- IntegrationSuite
- OpenSystemPlaced
- sync.Mutex
- net_tester/main.go
- render.go
- time.Time
- IntegrationSuite
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- NewEventFanout
- Read
- preload_wait.c
- validateConfigText
- runNetworkGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- liveSession
- NetworkMonitorUseCase
- IntegrationSuite
- Well-known bypass classes
- newTestVetter
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- runMonitor
- Build & sign release assets (reusable workflow)
- runGenKey
- integrationtests/guard_test.go
- Vault
- model
- app-listener Terminal Demo Recording
- expandPlaceholders
- vaultFS
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- RawTarget
- watchPattern
- WebSocket /ws client
- NewDumpHook
- writeWithin
- install.sh
- changelog.go
- launch_rule_at
- make test-integration (rootful Docker suite)
- netGuardModel
- supersede_probe.c
- fileedit_symlink_test.go
- trace-app-libs.sh
- IntegrationSuite
- raw_block_device.c
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- tempGrantSetup
- bottomBar
- find.go
- io.Reader
- Guard
- tempRow
- headerWidget
- configedit_test.go
- mustSubkey
- downloadAndVerify
- findTopLevelWindow
- ResolveCatalogEntry
- encryptDirectories
- find_test.go
- runGuard
- formatEntry
- MulticallError
- Vault
- graphify-refresh.sh
- TemporaryGrant
- ldConfParser
- ResolveTempBinary
- copyXattrs
- golang.org/x/sys/unix.Timespec

## God Nodes (most connected - your core abstractions)
1. `Guard` - 93 edges
2. `Config` - 91 edges
3. `Load()` - 85 edges
4. `GuardInodeKey` - 82 edges
5. `IntegrationSuite` - 75 edges
6. `absPath()` - 74 edges
7. `IntegrationSuite` - 66 edges
8. `writeConfig()` - 65 edges
9. `fileEditModel` - 60 edges
10. `User` - 56 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (172 total, 33 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (16): Execute(), tempJournalRow, launchRule, trustEvent, btrfsIoctlFsInfoArgs, Check(), collectDiagnostics(), readSysctl() (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (117): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+109 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.06
Nodes (13): IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe() (+5 more)

### Community 3 - "Vault"
Cohesion: 0.11
Nodes (21): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+13 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (72): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+64 more)

### Community 5 - "testKeyPair"
Cohesion: 0.29
Nodes (8): parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParsePublicKey(), TestVerifyRelease(), verifyRelease(), verifyReleaseWithKey()

### Community 6 - "monitor.bpf.c"
Cohesion: 0.08
Nodes (59): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+51 more)

### Community 7 - "User"
Cohesion: 0.12
Nodes (38): trustManager, addKeysToAgent(), installSSHAgentEnv(), removeSSHAgentEnv(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent() (+30 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (28): addErrKind, guardUnitTest, ReplacementCheck, BinaryEntry, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder() (+20 more)

### Community 9 - "mkdirs"
Cohesion: 0.14
Nodes (35): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+27 more)

### Community 10 - "Hash"
Cohesion: 0.11
Nodes (33): auditAfterEdit(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+25 more)

### Community 11 - "CandidateDir"
Cohesion: 0.19
Nodes (9): entryWriters(), reserveChain(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), CandidateDir, inOwnTree() (+1 more)

### Community 12 - "trustManager"
Cohesion: 0.11
Nodes (19): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), trustManager (+11 more)

### Community 13 - "runUninstall"
Cohesion: 0.09
Nodes (34): parseMaintenanceFlags(), runMaintenanceMode(), TestValidateMaintenanceFlags(), validateMaintenanceFlags(), deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall() (+26 more)

### Community 14 - "highlight_test.go"
Cohesion: 0.15
Nodes (16): Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned() (+8 more)

### Community 15 - "buildOneGuard"
Cohesion: 0.08
Nodes (39): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded() (+31 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "GenerateConf"
Cohesion: 0.05
Nodes (62): TestDiffMergePreservesExistingAndParses(), groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), updateCatalogConfig(), LibraryBlock, RefreshOptions (+54 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "New"
Cohesion: 0.17
Nodes (13): diffModel, editorModel, moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys() (+5 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.16
Nodes (45): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+37 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (11): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), uniqueEncryptionRoots() (+3 more)

### Community 23 - "KernelDev"
Cohesion: 0.09
Nodes (31): fdStat, rootHandle, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), openRoot(), physicalParent() (+23 more)

### Community 24 - "Candidate"
Cohesion: 0.09
Nodes (35): appendSectionsAndEdit(), collectDiffAdditions(), runDiffCatalog(), selectAndEditConfig(), cleanOrphanedFscrypt(), addManualDirectories(), editConfig(), groupCandidates() (+27 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.11
Nodes (48): add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector(), discover_guarded_parent() (+40 more)

### Community 26 - "Run"
Cohesion: 0.22
Nodes (6): columnState, dataRow, guiModel, newColumnState(), newDataRow(), Run()

### Community 27 - "NetEventType"
Cohesion: 0.07
Nodes (20): NetBpfEvent, NetEvent, NetEventType, foreignPinnedLink(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType() (+12 more)

### Community 28 - "SanitizeText"
Cohesion: 0.06
Nodes (30): newBinaryVetter(), openBinaryVetter(), ownerPaths(), logRefreshChange(), TestLogRefreshChange(), grantNewBinaries(), runHeadless(), runHeadless() (+22 more)

### Community 29 - "unlockUnderGuard"
Cohesion: 0.50
Nodes (5): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), SectionScan

### Community 30 - "SelectOrphans"
Cohesion: 0.24
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 31 - "serve.go"
Cohesion: 0.11
Nodes (22): basicAuth(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop(), relayServeSignals() (+14 more)

### Community 32 - "engine"
Cohesion: 0.08
Nodes (9): TestTieRank(), engine, pathWithin(), tieRank(), deleteInoKeys(), guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), intersectAllows() (+1 more)

### Community 33 - "codeedit.go"
Cohesion: 0.25
Nodes (17): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+9 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.15
Nodes (27): fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd() (+19 more)

### Community 35 - "DaemonEvent"
Cohesion: 0.07
Nodes (44): ServeConfig, UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), notifySystemdReady() (+36 more)

### Community 36 - "daemonconfig.go"
Cohesion: 0.16
Nodes (38): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+30 more)

### Community 37 - "update.go"
Cohesion: 0.09
Nodes (32): changelogText(), showChangelog(), TestChangelogText(), applyUpdate(), compareStableVersions(), confirmUpdate(), fetchReleases(), filterChannel() (+24 more)

### Community 38 - "GuardEvent"
Cohesion: 0.12
Nodes (10): bpfGuardEvent, GuardEvent, commMatchesGuardedBinary(), fsGateLabel(), logBacklogDenial(), parseGuardEvent(), processGateLabel(), TestParseGuardEventFsGateLabels() (+2 more)

### Community 39 - "guarded_ancestor_within_limit"
Cohesion: 0.26
Nodes (17): evict_inode_from_guard(), get_dentry_from_path(), get_inode_from_path(), guard_inode_getattr(), guard_path_link(), guard_path_mkdir(), guard_path_mknod(), guard_path_rename() (+9 more)

### Community 40 - "sanitizeTerminalText"
Cohesion: 0.17
Nodes (8): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), TestSanitizeTerminalText()

### Community 41 - ".GrantTemporaryAccess"
Cohesion: 0.09
Nodes (14): fakeDaemon, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess, logTempGrant(), newExeTap(), planTempGrants() (+6 more)

### Community 44 - "main_test.go"
Cohesion: 0.06
Nodes (34): decodeMulticall(), IntegrationSuite, multicallConfig(), uutilsLinks(), guardBinaryFlag(), infraContainerPaths(), tailLast(), TestIntegrationSuite() (+26 more)

### Community 45 - "os.FileMode"
Cohesion: 0.10
Nodes (29): TestDaemonUnitWithMetadataOutput(), daemonUnitWithMetadataOutput(), dupFD(), installConfig(), installFile(), installFileAs(), installServices(), installSSHAgent() (+21 more)

### Community 46 - "Monitor"
Cohesion: 0.06
Nodes (8): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), Monitor, monitorUnitTest, pathCache, pathCacheEntry

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 50 - "filevault.go"
Cohesion: 0.21
Nodes (14): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+6 more)

### Community 51 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 52 - "IntegrationSuite"
Cohesion: 0.11
Nodes (5): IntegrationSuite, exploitTest, IntegrationSuite, newEventTypes(), IntegrationSuite

### Community 53 - "runInstall"
Cohesion: 0.12
Nodes (20): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), configPaths(), pathCovered(), TestPathCovered() (+12 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.12
Nodes (7): clamp(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel, servedModel, serveTestModel, sessionModel

### Community 55 - "absPath"
Cohesion: 0.10
Nodes (5): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath()

### Community 56 - "os.File"
Cohesion: 0.07
Nodes (42): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+34 more)

### Community 57 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 58 - "runForward"
Cohesion: 0.16
Nodes (15): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+7 more)

### Community 59 - "EventType"
Cohesion: 0.09
Nodes (11): BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryEventsReadOnlyRejectsRestriction(), Run(), Cstr(), DecodeBpfEvent() (+3 more)

### Community 60 - "TrustGuard"
Cohesion: 0.07
Nodes (17): GlobChild, GlobReservations, statFunc, trustHook, trustRows, UpdaterPlan, setSupersedeTrust(), cStr() (+9 more)

### Community 61 - "GuardInodeKey"
Cohesion: 0.05
Nodes (39): Inspector, refusedReplacements, SupersededBinary, checkBtrfsKeys(), ensureBtrfsLayout(), engine, TrustGuard, SetInspectors() (+31 more)

### Community 62 - "inode_dev"
Cohesion: 0.12
Nodes (24): ctx_ptr(), inode_dev(), sb_dev(), exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused() (+16 more)

### Community 63 - "fakeGuardRepo"
Cohesion: 0.08
Nodes (8): multicallAdmissible(), systemPatterns(), systemRule(), vetDecision, Line(), BinaryRule, ParseEventType(), fakeGuardRepo

### Community 64 - "string"
Cohesion: 0.14
Nodes (4): denied(), main(), denied(), main()

### Community 65 - "pinstate.go"
Cohesion: 0.07
Nodes (44): CheckBPFLSM(), configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), newPinGeneration(), prepareDaemonStart(), relockStaleVaults(), runBPFCheck() (+36 more)

### Community 66 - "LibraryClosure"
Cohesion: 0.12
Nodes (17): defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldSoConfDirs(), LdSoPreloadPaths(), LibraryClosure(), resolveSoname() (+9 more)

### Community 67 - "time.Duration"
Cohesion: 0.21
Nodes (4): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), debouncer

### Community 68 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 69 - "catalogRefresher"
Cohesion: 0.21
Nodes (5): drain(), newCatalogRefresher(), readEvents(), sameConfig(), catalogRefresher

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 72 - "bufio.Reader"
Cohesion: 0.15
Nodes (13): swapFile(), controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines() (+5 more)

### Community 73 - "auditGroupCoverage"
Cohesion: 0.20
Nodes (12): auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+4 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "filevault_test.go"
Cohesion: 0.22
Nodes (22): looksLikeFileVaultRecord(), inode(), TestEnsureRecoverySidecarPlaceholder(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestStageRecoverySealsContentAndStaysInPlace(), TestUnlockFileInPlaceRecoversFromInterruptedTransform() (+14 more)

### Community 76 - "serveWebSocket"
Cohesion: 0.21
Nodes (7): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), snapshotHub

### Community 78 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (10): copy(), is_verb(), main(), usage(), dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd() (+2 more)

### Community 80 - "runOfflineEdit"
Cohesion: 0.18
Nodes (11): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), decryptStep(), pickDirsToDecrypt(), VerifyEncryptedKeys(), ScanEncryptedCatalogDirs() (+3 more)

### Community 81 - "systemd.go"
Cohesion: 0.10
Nodes (37): applyDiffAdditions(), restoreDaemonAfterDiffAbort(), applyLiveRefresh(), buildBinaryIfNeeded(), installBinaryOnly(), runUpdateCatalogOnly(), softenAutomatedRefreshErr(), TestApplyLiveRefreshDeliversReload() (+29 more)

### Community 83 - "OpenSystemPlaced"
Cohesion: 0.09
Nodes (21): setConfinedHomes(), placedWalk, TrustGuard, launchKey(), TestLaunchRulesFitTheKey(), confineBelow(), inUserHome(), SetUserHomes() (+13 more)

### Community 84 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 85 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 86 - "render.go"
Cohesion: 0.13
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "time.Time"
Cohesion: 0.21
Nodes (8): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, TestNearestRootClaims(), nearestRootClaims()

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "NewEventFanout"
Cohesion: 0.24
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 92 - "Read"
Cohesion: 0.10
Nodes (20): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+12 more)

### Community 94 - "validateConfigText"
Cohesion: 0.09
Nodes (23): collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+15 more)

### Community 95 - "runNetworkGuard"
Cohesion: 0.11
Nodes (18): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), DiscoverInfraBinaries(), Mode (+10 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "liveSession"
Cohesion: 0.24
Nodes (6): parseErr(), holdForward(), pingInterval(), liveSession, RunFileEditorSession(), RunSession()

### Community 100 - "NetworkMonitorUseCase"
Cohesion: 0.12
Nodes (13): NetworkMonitorRepository, NewGuardUseCase(), NetworkMonitorUseCase, NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+5 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "newTestVetter"
Cohesion: 0.14
Nodes (20): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSwapWithinGenerationLosesRights() (+12 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "runMonitor"
Cohesion: 0.11
Nodes (15): CheckEBPF(), ParseNetEventsFlag(), runHeadless(), runMonitor(), runTUI(), runNetworkMonitor(), runTUI(), MonitorRepository (+7 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runGenKey"
Cohesion: 0.22
Nodes (10): confirmMasterKeyOverwrite(), runGenKey(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen(), GenerateMasterKey() (+2 more)

### Community 110 - "integrationtests/guard_test.go"
Cohesion: 0.21
Nodes (6): attrOpCase, eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 111 - "Vault"
Cohesion: 0.18
Nodes (10): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), CopyTreeWithProgress(), TestCopyTreeProgress() (+2 more)

### Community 112 - "model"
Cohesion: 0.26
Nodes (6): formatType(), listenForEvents(), tickStats(), eventLine, eventMsg, model

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 114 - "expandPlaceholders"
Cohesion: 0.10
Nodes (16): symlinkTargetGlob(), confSafeMatch(), expandPlaceholders(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), CandidateDir, ConfSafePath() (+8 more)

### Community 115 - "vaultFS"
Cohesion: 0.13
Nodes (14): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode(), listEntries(), openRootT(), TestDefaultNewFileMode() (+6 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 121 - "watchPattern"
Cohesion: 0.26
Nodes (9): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+1 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 124 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "changelog.go"
Cohesion: 0.23
Nodes (6): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 127 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "netGuardModel"
Cohesion: 0.29
Nodes (3): protoString(), netGuardEventMsg, netGuardModel

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "raw_block_device.c"
Cohesion: 0.83
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "tempGrantSetup"
Cohesion: 0.33
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 143 - "bottomBar"
Cohesion: 0.23
Nodes (7): decryptDirectories(), TestDecryptDirectoriesEmpty(), TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 144 - "find.go"
Cohesion: 0.27
Nodes (7): findButton, span, findAll(), lowerRunes(), newFindInput(), runesEqual(), TestFindAllSmartCase()

### Community 145 - "io.Reader"
Cohesion: 0.25
Nodes (6): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype(), progressReader

### Community 146 - "Guard"
Cohesion: 0.06
Nodes (20): binaryVerifyState, deferredBinary, Mode, sharedHashEntry, TestWithChmodDropWriteNeedsReadOnly(), SharedPinDegraded(), canonicalBinaryPath(), canonicalPaths() (+12 more)

### Community 147 - "tempRow"
Cohesion: 0.17
Nodes (8): exeRow, tempGrant, tempMaskOp, tempRow, TempRule, Guard, liveTrust(), planTempBlock()

### Community 149 - "configedit_test.go"
Cohesion: 0.57
Nodes (7): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession()

### Community 150 - "mustSubkey"
Cohesion: 0.31
Nodes (9): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat() (+1 more)

### Community 151 - "downloadAndVerify"
Cohesion: 0.31
Nodes (7): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), sanityCheckBinary(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), updateHTTPClient()

### Community 152 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 153 - "ResolveCatalogEntry"
Cohesion: 0.60
Nodes (5): catalogEntryMatch(), findCatalogRoot(), findCatalogWatchSubPath(), isInsidePath(), ResolveCatalogEntry()

### Community 154 - "encryptDirectories"
Cohesion: 0.38
Nodes (4): encryptDirectories(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 155 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 156 - "runGuard"
Cohesion: 0.11
Nodes (16): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), displayEntries(), holdHeadless(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule() (+8 more)

### Community 157 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 165 - "ResolveTempBinary"
Cohesion: 0.09
Nodes (18): TempBinary, BinaryStat, deleteResKeys(), syncU32Map(), binaryStatOf(), loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow() (+10 more)

### Community 166 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

## Knowledge Gaps
- **53 isolated node(s):** `guardExploitTest`, `attrOpCase`, `TrustGuard`, `engine`, `Guard` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 391 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **33 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`, `integrationtests/guard_test.go`, `.startGuardStd`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `DaemonEvent`, `GuardEvent`, `.GrantTemporaryAccess`, `buildOneGuard`, `GuardInodeKey`, `sync.Mutex`, `time.Time`, `os.File`, `unlockUnderGuard`, `fakeGuardRepo`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Why does `shQuote()` connect `IntegrationSuite` to `daemon/daemon.go`, `IntegrationSuite`, `main_test.go`, `IntegrationSuite`, `absPath`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **What connects `guardExploitTest`, `attrOpCase`, `TrustGuard` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10274960837661802 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.025070028011204483 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.06158415841584158 - nodes in this community are weakly interconnected._