# Graph Report - app-listener  (2026-10-12)

## Corpus Check
- 359 files · ~1,558,243 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4657 nodes · 18361 edges · 201 communities (160 shown, 41 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1341 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `295782b7`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- testcontainers.Container
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- setpassword.go
- planUpdaters
- TestMonitorUseCaseLifecycle
- Candidate
- install.go
- Vault
- IntegrationSuite
- guard_trust.bpf.c
- confgen.go
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- eventMask
- buildOneGuard
- guard.bpf.c
- NetworkMonitor
- find_test.go
- Ledger
- filevault_test.go
- Monitor
- Serve
- engine
- IntegrationSuite
- fileedit_model_test.go
- IntegrationSuite
- Run
- uninstall.go
- controlServer
- bpfstats/main.go
- Finder
- IntegrationSuite
- IntegrationSuite
- liveSession
- IntegrationSuite
- update.go
- highlight_test.go
- .syncSetsLocked
- TrustGuard
- tempRow
- inode_dev
- IntegrationSuite
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- absPath
- codeedit.go
- sync.Mutex
- mustSubkey
- IntegrationSuite
- auditGroupCoverage
- runNetworkGuard
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- migrate.go
- CandidateDir
- filevault.go
- time.Duration
- fcntl
- copytree.go
- trustManager
- uninstall_system.go
- btrfsLayoutFrom
- monitorUnitTest
- github.com/spf13/cobra.Command
- runNetworkMonitor
- process_vm_readv.c
- golang.org/x/sys/unix.Stat_t
- .resyncLink
- Read
- Vault
- MountVouch
- lockAndDeprovision
- github.com/charmbracelet/bubbles/textarea.Model
- time.Time
- fakeDaemon
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- os.FileMode
- KernelDev
- IntegrationSuite
- Well-known bypass classes
- runGuard
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- find.go
- Build & sign release assets (reusable workflow)
- runMonitor
- fileedit_symlink_test.go
- net_tester/main.go
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- .pathInfo
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- catalogRefresher
- runDaemon
- WebSocket /ws client
- netModel
- string
- install.sh
- mc_tag.h
- changelog.go
- make test-integration (rootful Docker suite)
- EventType
- supersede_probe.c
- bufio.Reader
- trace-app-libs.sh
- IntegrationSuite
- walkInodes
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- ClassifyMulticall
- os.File
- preload_wait.c
- Guard
- writeWithin
- launchKey
- fscrypt/orphans.go
- fscrypt.go
- validateFlags
- daemon: permanent protection with encryption at rest
- newTestVetter
- GuardInodeKey
- unlockUnderGuard
- New
- testcontainers.Container
- Installation
- ExeKey
- graphify-refresh.sh
- NewGuard
- RefreshCatalog
- ConfinedEntry
- GuardEvent
- .runEditorHarness
- loadDaemonConfig
- inodecap_test.go
- NetGuard
- Limitations
- Troubleshooting
- inodeChain
- integrationtests/guard_test.go
- Daemon configuration
- edit-protected
- NewDumpHook
- Compatibility
- How it works
- harnessSuite
- inodeRowsOf
- bpfLsmListed
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS

## God Nodes (most connected - your core abstractions)
1. `Guard` - 107 edges
2. `Config` - 95 edges
3. `Load()` - 89 edges
4. `IntegrationSuite` - 81 edges
5. `absPath()` - 79 edges
6. `writeConfig()` - 69 edges
7. `IntegrationSuite` - 66 edges
8. `fileEditModel` - 60 edges
9. `User` - 56 edges
10. `GuardInodeKey` - 54 edges

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

## Communities (201 total, 41 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.10
Nodes (26): Execute(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey (+18 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (115): plantBun(), TestOpenBunTmp_KeepsVettedDir(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait() (+107 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.12
Nodes (7): IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), shQuote(), daemonEvent

### Community 3 - "testcontainers.Container"
Cohesion: 0.16
Nodes (3): configBinaryLines(), rawExec(), watchPathsInConfig()

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (78): TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines() (+70 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.14
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.07
Nodes (68): trustManager, askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort() (+60 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.07
Nodes (23): guardUnitTest, ReplacementCheck, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), copySelf(), runTool(), TestReadOnlyTreeInodeDoesNotFollowSymlinks() (+15 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 10 - "setpassword.go"
Cohesion: 0.09
Nodes (41): auditAfterEdit(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+33 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+4 more)

### Community 13 - "Candidate"
Cohesion: 0.13
Nodes (25): appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates() (+17 more)

### Community 14 - "install.go"
Cohesion: 0.07
Nodes (54): applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly() (+46 more)

### Community 15 - "Vault"
Cohesion: 0.19
Nodes (10): isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+2 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "confgen.go"
Cohesion: 0.05
Nodes (61): logRefreshChange(), TestLogRefreshChange(), TestDiffMergePreservesExistingAndParses(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH() (+53 more)

### Community 19 - "fileEditModel"
Cohesion: 0.11
Nodes (10): clipLabel(), humanSize(), isBinaryContent(), sortNodes(), TestIsBinaryContent(), chownUser, fileEditMode, fileEditModel (+2 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.26
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.11
Nodes (54): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+46 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (13): ownerPaths(), Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots() (+5 more)

### Community 23 - "eventMask"
Cohesion: 0.09
Nodes (16): deferredBinary, VettedInode, BinariesSummary(), canonicalBinaryPath(), confinedEntry(), eventMask(), patternMatches(), readPinnedResID() (+8 more)

### Community 24 - "buildOneGuard"
Cohesion: 0.10
Nodes (29): buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), grantNewBinaries() (+21 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (53): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class() (+45 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.11
Nodes (10): NetBpfEvent, NetEvent, foreignPinnedLink(), FormatAddr(), Ntohs(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure() (+2 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "Ledger"
Cohesion: 0.06
Nodes (27): newBinaryVetter(), openBinaryVetter(), candidates(), confirmable(), confirmAll(), describe(), hashNow(), run() (+19 more)

### Community 29 - "filevault_test.go"
Cohesion: 0.30
Nodes (17): looksLikeFileVaultRecord(), inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsFileVaultRecordShape(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceRecoversFromInterruptedTransform(), TestUnlockFileInPlaceWrongKeyFails() (+9 more)

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (28): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+20 more)

### Community 32 - "engine"
Cohesion: 0.08
Nodes (10): TestNearestRootClaims(), TestTieRank(), ConfigureSharedPinning(), engine, nearestRootClaims(), pathWithin(), TestPathWithinMatchesPrefixRule(), tieRank() (+2 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.16
Nodes (25): fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod() (+17 more)

### Community 35 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 36 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 37 - "uninstall.go"
Cohesion: 0.10
Nodes (31): deletePostBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), offerBackupCleanup(), pickDirsToDecrypt() (+23 more)

### Community 38 - "controlServer"
Cohesion: 0.20
Nodes (7): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, forwardRequest, grantRequest

### Community 39 - "bpfstats/main.go"
Cohesion: 0.12
Nodes (24): confirmMasterKeyOverwrite(), runBPFCheck(), runGenKey(), runOneShotMode(), guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad() (+16 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.12
Nodes (9): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent() (+1 more)

### Community 45 - "update.go"
Cohesion: 0.06
Nodes (52): applyUpdate(), assetsFor(), checkUpdatePreconditions(), compareStableVersions(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), fetchReleases() (+44 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.14
Nodes (17): fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo(), TestDaemonConfTokens() (+9 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.11
Nodes (18): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, intersectAllows(), isSubset(), planMemberRows() (+10 more)

### Community 49 - "TrustGuard"
Cohesion: 0.15
Nodes (5): trustRows, setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 50 - "tempRow"
Cohesion: 0.17
Nodes (7): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust(), planTempBlock()

### Community 51 - "inode_dev"
Cohesion: 0.14
Nodes (19): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), code_suspect(), current_is_inspector() (+11 more)

### Community 52 - "IntegrationSuite"
Cohesion: 0.06
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.15
Nodes (39): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+31 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.09
Nodes (11): clamp(), listenForGuardEvents(), tickGuardStats(), renderBullets(), TestRenderBulletsVerbatimAndIndented(), guardEventMsg, netGuardEventMsg, noticeModel (+3 more)

### Community 55 - "absPath"
Cohesion: 0.09
Nodes (8): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath(), TestMain(), helperChild(), TestMain()

### Community 56 - "codeedit.go"
Cohesion: 0.25
Nodes (17): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), Navigate() (+9 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 58 - "mustSubkey"
Cohesion: 0.33
Nodes (8): newFileVaultAEAD(), openFileVault(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestOpenFileVaultRejectsTampering(), TestSealFileVaultNoncesNeverRepeat(), TestSealOpenFileVaultRoundTrip()

### Community 60 - "auditGroupCoverage"
Cohesion: 0.20
Nodes (12): auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+4 more)

### Community 61 - "runNetworkGuard"
Cohesion: 0.12
Nodes (16): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode, NetGuardEvent (+8 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.18
Nodes (17): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+9 more)

### Community 63 - "fakeGuardRepo"
Cohesion: 0.07
Nodes (9): dropSupersededLocked(), Guard, SupersededBinary, PruneSuperseded(), supersededMark(), supersedeMaps(), supersedePruneLoop(), warnLostStamps() (+1 more)

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.14
Nodes (26): configureDaemonLogging(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState() (+18 more)

### Community 66 - "migrate.go"
Cohesion: 0.35
Nodes (4): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), hint

### Community 67 - "CandidateDir"
Cohesion: 0.06
Nodes (29): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName(), confSafeMatch() (+21 more)

### Community 68 - "filevault.go"
Cohesion: 0.21
Nodes (16): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace() (+8 more)

### Community 69 - "time.Duration"
Cohesion: 0.13
Nodes (9): debouncer, New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub (+1 more)

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copytree.go"
Cohesion: 0.20
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "trustManager"
Cohesion: 0.13
Nodes (20): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), trustManager (+12 more)

### Community 73 - "uninstall_system.go"
Cohesion: 0.29
Nodes (11): detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), removeSSHAgentUnit(), revertSSHAgents(), staleRCUsers() (+3 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.12
Nodes (23): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+15 more)

### Community 78 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (9): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+1 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.15
Nodes (15): BinaryStat, placedWalk, binaryStatOf(), binaryStatOf(), StatBinary(), StatBinaryFile(), inUserHome(), GlobSystemPlaced() (+7 more)

### Community 81 - ".resyncLink"
Cohesion: 0.16
Nodes (9): binaryVerifyState, sharedHashEntry, hashFileShared(), hashPathShared(), settled(), unchangedSince(), engine, Guard (+1 more)

### Community 82 - "Read"
Cohesion: 0.06
Nodes (31): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+23 more)

### Community 83 - "Vault"
Cohesion: 0.22
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 84 - "MountVouch"
Cohesion: 0.14
Nodes (8): deleteInoKeys(), deleteResKeys(), syncU32Map(), syncMap(), deleteKeys(), MountVouch, TrustGuard, NewMountVouch()

### Community 85 - "lockAndDeprovision"
Cohesion: 0.21
Nodes (8): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), modifiedContextWithSource()

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.17
Nodes (13): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 87 - "time.Time"
Cohesion: 0.12
Nodes (16): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+8 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.11
Nodes (31): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+23 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.11
Nodes (11): NetEventType, FirstOnBtrfs(), NetEventTypes(), ParseNetEventType(), eventsetSummary(), eventTypeKey(), NetGuard, mcNameKey() (+3 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "os.FileMode"
Cohesion: 0.09
Nodes (30): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig(), TestDaemonUnitWithMetadataOutput() (+22 more)

### Community 100 - "KernelDev"
Cohesion: 0.08
Nodes (32): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+24 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "runGuard"
Cohesion: 0.10
Nodes (22): CheckBPFLSM(), CheckEBPF(), newPinGeneration(), prepareDaemonStart(), displayEntries(), holdForward(), holdHeadless(), HeadlessLine() (+14 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "find.go"
Cohesion: 0.20
Nodes (7): findButton, findAll(), lowerRunes(), newFindInput(), overlay(), runesEqual(), TestFindAllSmartCase()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.09
Nodes (22): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), pathWithinDir(), ResolveTargets(), validatePair() (+14 more)

### Community 110 - "fileedit_symlink_test.go"
Cohesion: 0.17
Nodes (18): TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), openRootT(), TestWriteFileKeepMeta(), TestWriteFileKeepMetaRefusesSymlink(), mustRead(), mustWrite() (+10 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NewEventFanout"
Cohesion: 0.25
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout[T]

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - ".pathInfo"
Cohesion: 0.29
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "catalogRefresher"
Cohesion: 0.18
Nodes (6): drain(), newCatalogRefresher(), readEvents(), sameConfig(), swapFile(), catalogRefresher

### Community 121 - "runDaemon"
Cohesion: 0.08
Nodes (34): backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), lockOneRoot(), lockRootRecovering(), makeReloadHandler(), noteReloadedSettings() (+26 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "netModel"
Cohesion: 0.21
Nodes (8): formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents(), tickNetStats(), netEventLine, netEventMsg, netModel

### Community 124 - "string"
Cohesion: 0.16
Nodes (6): main(), report(), copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 127 - "changelog.go"
Cohesion: 0.18
Nodes (11): changelogText(), newChangelogModel(), showChangelog(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestNewChangelogModelViewContainsTitleAndNotes() (+3 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.09
Nodes (12): BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), Run(), Cstr() (+4 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "bufio.Reader"
Cohesion: 0.19
Nodes (11): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+3 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "walkInodes"
Cohesion: 0.15
Nodes (7): addErrKind, classifyAddErr(), addRecorder(), walkEntries(), walkInodes(), walkLiveEntries(), TestReadOnlyWalkToleratesDanglingSymlinks()

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.10
Nodes (21): confirmUserPlaced(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps(), ResolveTempBinary(), stripPinnedAllow(), StripPinnedTempAllows() (+13 more)

### Community 143 - "ClassifyMulticall"
Cohesion: 0.10
Nodes (21): GlobChild, statFunc, Multicall, MulticallKind, multicallFamily(), singleBinaryPaths(), compileGlobNames(), globText() (+13 more)

### Community 144 - "os.File"
Cohesion: 0.11
Nodes (23): multicallAdmissible(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched() (+15 more)

### Community 146 - "Guard"
Cohesion: 0.09
Nodes (6): openRoot(), SharedPinDegraded(), Guard, mountinfoDev(), SuperblockDev(), TestMountinfoDev()

### Community 147 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 148 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 149 - "fscrypt/orphans.go"
Cohesion: 0.23
Nodes (19): appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 150 - "fscrypt.go"
Cohesion: 0.06
Nodes (40): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+32 more)

### Community 151 - "validateFlags"
Cohesion: 0.40
Nodes (5): TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "newTestVetter"
Cohesion: 0.33
Nodes (10): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+2 more)

### Community 154 - "GuardInodeKey"
Cohesion: 0.09
Nodes (16): refusedReplacements, rootHandle, engine, Inspector, SetInspectors(), appletTag(), buildOf(), engine (+8 more)

### Community 155 - "unlockUnderGuard"
Cohesion: 0.50
Nodes (5): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), SectionScan

### Community 156 - "New"
Cohesion: 0.17
Nodes (13): diffModel, editorModel, moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys() (+5 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "ExeKey"
Cohesion: 0.08
Nodes (31): setConfinedHomes(), TrustGuard, ExeKey(), IsMulticallRefusal(), MulticallRefusal(), RefuseMulticallLines(), registerBuild(), TestExeKeyRefusesGNUCoreutilsSingle() (+23 more)

### Community 162 - "NewGuard"
Cohesion: 0.22
Nodes (7): Mode, TestWithChmodDropWriteNeedsReadOnly(), canonicalPaths(), guardModeKey(), modeString(), NewGuard(), WithChmodDropWrite()

### Community 163 - "RefreshCatalog"
Cohesion: 0.40
Nodes (5): RefreshCatalog(), TestRefreshCatalogAdmitVetsOnlyNewLines(), TestRefreshCatalogKeepsSymlinkLineNotTarget(), TestRefreshCatalogReportsVersionBump(), TestRefreshCatalogUserSectionsOnly()

### Community 164 - "ConfinedEntry"
Cohesion: 0.42
Nodes (11): ConfinedEntry(), hashedReused(), resetSharedHashes(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestSharedHashMissesOnChangedFingerprint(), TestSharedHashNotPublishedBeforeSettled(), TestSharedHashNotPublishedWhenWrittenDuringRead() (+3 more)

### Community 165 - "GuardEvent"
Cohesion: 0.13
Nodes (12): bpfGuardEvent, GuardEvent, commMatchesGuardedBinary(), decodeGuardEvent(), fsGateLabel(), logBacklogDenial(), parseGuardEvent(), processGateLabel() (+4 more)

### Community 167 - "loadDaemonConfig"
Cohesion: 0.67
Nodes (3): loadDaemonConfig(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup()

### Community 169 - "NetGuard"
Cohesion: 0.33
Nodes (3): trustHook, NetGuard, hook

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 175 - "integrationtests/guard_test.go"
Cohesion: 0.31
Nodes (6): attrOpCase, eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "NewDumpHook"
Cohesion: 0.33
Nodes (6): BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity(), TestDumpHookFormat()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 183 - "inodeRowsOf"
Cohesion: 0.15
Nodes (14): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+6 more)

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 186 - "Binary trust"
Cohesion: 0.33
Nodes (6): Binary ledger, Binary replacement, Binary trust, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.27
Nodes (5): applyNewFileMeta(), defaultNewFileMode(), listEntries(), TestDefaultNewFileMode(), vaultFS

## Knowledge Gaps
- **146 isolated node(s):** `Check your host first`, `Distributions`, `BPF-LSM must be active, not just compiled in`, `Mandatory and best-effort hooks`, `Architectures` (+141 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 501 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **41 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `shQuote()` connect `IntegrationSuite` to `daemon/daemon.go`, `testcontainers.Container`, `.runEditorHarness`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `absPath`, `IntegrationSuite`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `engine`, `daemon/daemon.go`, `NewGuard`, `time.Duration`, `GuardEvent`, `fakeDaemon`, `os.File`, `.resyncLink`, `time.Time`, `eventMask`, `buildOneGuard`, `runDaemon`, `unlockUnderGuard`, `DaemonEvent`, `sync.Mutex`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `testcontainers.Container`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `testcontainers.Container`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **What connects `Check your host first`, `Distributions`, `BPF-LSM must be active, not just compiled in` to the rest of the system?**
  _146 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1016966142896524 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.025067654180316194 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.12439024390243902 - nodes in this community are weakly interconnected._