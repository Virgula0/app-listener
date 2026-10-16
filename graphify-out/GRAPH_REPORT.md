# Graph Report - app-listener  (2026-10-16)

## Corpus Check
- 379 files · ~1,571,023 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4754 nodes · 19958 edges · 194 communities (162 shown, 32 thin omitted)
- Extraction: 92% EXTRACTED · 8% INFERRED · 0% AMBIGUOUS · INFERRED: 1667 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `a64be1b6`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- GuardInodeKey
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- User
- guardUnitTest
- globreserve_test.go
- Guard
- liveSession
- go_pkg_path_filepath
- Config
- install.go
- Candidate
- .startGuardStd
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- IntegrationSuite
- guard.bpf.c
- NetworkMonitor
- highlight_test.go
- LibraryClosure
- main_test.go
- Vault
- serve.go
- bpfstats/main.go
- changelog.go
- fileedit_model_test.go
- Ledger
- uninstall_system.go
- sync.Mutex
- ConfinedEntry
- vaultFS
- GenerateConf
- IntegrationSuite
- monitorUnitTest
- newWatchPattern
- update_test.go
- netGuardModel
- stdio
- time.Time
- .FilterExistingWhitelist
- planUpdaters
- codeedit.go
- model
- watchGroup
- github.com/charmbracelet/bubbletea.Model
- TrustGuard
- system.go
- absPath
- engine
- github.com/charmbracelet/bubbles/textarea.Model
- ResolveTempBinary
- fileedit_symlink_test.go
- exe_supersede.h
- GuardEvent
- expandPlaceholders
- pinstate.go
- gui/session.go
- fakeGuardRepo
- rootHandle
- buildTrustedSet
- fcntl
- copyRegularAt
- KernelDev
- globBuilder
- btrfsLayoutFrom
- SetSectionWhitelist
- github.com/spf13/cobra.Command
- IntegrationSuite
- highlight.go
- process_vm_readv.c
- net_tester/main.go
- hashCache
- CandidateDir
- runUninstall
- IntegrationSuite
- netModel
- audit.go
- TempBinary
- .refreshOnce
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- time.Duration
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- github.com/charmbracelet/bubbletea.Cmd
- golang.org/x/sys/unix.Stat_t
- IntegrationSuite
- Well-known bypass classes
- Pulse
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- NetGuardEvent
- Build & sign release assets (reusable workflow)
- runMonitor
- github.com/stretchr/testify/suite.Suite
- controlManager
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- eventUnitTest
- hashFileShared
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- compileGlobNames
- NewDumpHook
- WebSocket /ws client
- TestBatchInodeRows
- launch_probe.c
- install.sh
- mc_tag.h
- fakeVault
- make test-integration (rootful Docker suite)
- writeEvent
- supersede_probe.c
- loadDaemonConfig
- trace-app-libs.sh
- IntegrationSuite
- writeWithin
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- testLimiter
- fakeDaemon
- mountTable
- buildOneGuard
- IntegrationSuite
- render.go
- bottomBar
- newTestVetter
- filevault.go
- Monitor
- daemon: permanent protection with encryption at rest
- .runAttrOp
- checkBtrfsKeys
- io.Reader
- find_test.go
- StripPinnedTempAllows
- Installation
- os.File
- graphify-refresh.sh
- RawTarget
- configedit_test.go
- runNetworkGuard
- inodeChain
- bpfLsmListed
- Vault
- BinaryStat
- Limitations
- Troubleshooting
- Execute
- mountsUnitTest
- Daemon configuration
- edit-protected
- EventType
- Compatibility
- How it works
- openBunTmp
- pathCache
- SanitizeText
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .eventTypeAllowed
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- monitor: see what touches a path
- kernel-prof.sh

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 95 edges
3. `GuardInodeKey` - 94 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 67 edges
9. `fileEditModel` - 60 edges
10. `engine` - 58 edges

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

## Communities (194 total, 32 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.13
Nodes (25): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey, attrOpCase (+17 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (86): awaitStartupOrSignal(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup(), TestAwaitStartupOrSignalCompletes(), TestDiffCatalogFlagRegistered(), classifySupportError() (+78 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "GuardInodeKey"
Cohesion: 0.08
Nodes (25): refusedReplacements, appendInodeRow(), batchInodeRows(), inodeRowsOf(), nearestRootClaims(), tieRank(), newVerifyState(), supersedeUnnamed() (+17 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (77): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+69 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "User"
Cohesion: 0.07
Nodes (63): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir() (+55 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.07
Nodes (21): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf() (+13 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.25
Nodes (25): buildGlobReservations(), bitOf(), discordLibConfig(), hasChild(), mkdirs(), steamLibDirConfig(), TestBuildGlobReservations_BunRootsListed(), TestBuildGlobReservations_BunTmpdir() (+17 more)

### Community 10 - "Guard"
Cohesion: 0.05
Nodes (23): deferredBinary, BinaryEntry, TestCheckBinaryEventsReadOnlyRejectsRestriction(), TestWithChmodDropWriteNeedsReadOnly(), SharedPinDegraded(), canonicalBinaryPath(), canonicalPaths(), commMatchesGuardedBinary() (+15 more)

### Community 11 - "liveSession"
Cohesion: 0.18
Nodes (12): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), confirmUserPlaced(), displayEntries(), eventsLabel(), forwardBinaries() (+4 more)

### Community 12 - "go_pkg_path_filepath"
Cohesion: 0.27
Nodes (5): nextConfig(), isAbsoluteCandidate(), TestCatalogSanity(), EditText(), hint

### Community 13 - "Config"
Cohesion: 0.11
Nodes (31): ServeConfig, newCatalogRefresher(), sameConfig(), catchLifecycleSignals(), makeReloadHandler(), markReady(), newDaemonModel(), noteReloadedSettings() (+23 more)

### Community 14 - "install.go"
Cohesion: 0.06
Nodes (56): runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), encryptDirectories(), ensureInstalledBinary(), ensureMasterKey(), installBinaryOnly(), parseMaintenanceFlags() (+48 more)

### Community 15 - "Candidate"
Cohesion: 0.18
Nodes (19): appendSectionsAndEdit(), configPaths(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates() (+11 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (8): clipLabel(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind, fileNode

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (43): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+35 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (12): ownerPaths(), Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles() (+4 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (30): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+22 more)

### Community 24 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (70): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+62 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.07
Nodes (16): NetBpfEvent, NetEvent, NetEventType, FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType(), eventsetSummary() (+8 more)

### Community 27 - "highlight_test.go"
Cohesion: 0.18
Nodes (23): moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation(), TestEnterPastNinetyNineLines() (+15 more)

### Community 28 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 29 - "main_test.go"
Cohesion: 0.05
Nodes (20): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+12 more)

### Community 30 - "Vault"
Cohesion: 0.09
Nodes (31): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deprovisionKind (+23 more)

### Community 31 - "serve.go"
Cohesion: 0.09
Nodes (30): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+22 more)

### Community 32 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 33 - "changelog.go"
Cohesion: 0.14
Nodes (29): changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), showChangelog(), rel(), tags() (+21 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (26): OpenRegularNoFollow(), fileModeToUnixMode(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd() (+18 more)

### Community 35 - "Ledger"
Cohesion: 0.12
Nodes (10): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), Ledger, Pending (+2 more)

### Community 36 - "uninstall_system.go"
Cohesion: 0.05
Nodes (68): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt() (+60 more)

### Community 38 - "sync.Mutex"
Cohesion: 0.11
Nodes (20): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines(), startControlServer() (+12 more)

### Community 39 - "ConfinedEntry"
Cohesion: 0.26
Nodes (20): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+12 more)

### Community 40 - "vaultFS"
Cohesion: 0.12
Nodes (16): checkFreshDir(), openDirNoFollow(), renewDir(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode(), isBinaryContent(), listEntries() (+8 more)

### Community 41 - "GenerateConf"
Cohesion: 0.08
Nodes (34): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic() (+26 more)

### Community 44 - "newWatchPattern"
Cohesion: 0.18
Nodes (18): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+10 more)

### Community 45 - "update_test.go"
Cohesion: 0.06
Nodes (36): assetsFor(), checkUpdatePreconditions(), downloadAndVerify(), downloadFile(), downloadReleaseFiles(), fetchReleases(), parseChecksum(), parsePublicKey() (+28 more)

### Community 46 - "netGuardModel"
Cohesion: 0.11
Nodes (11): formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), formatGuardType(), clipToWidth(), formatResourceBar(), protoString() (+3 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "time.Time"
Cohesion: 0.16
Nodes (9): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID, sweepMemo (+1 more)

### Community 49 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 50 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 51 - "codeedit.go"
Cohesion: 0.25
Nodes (16): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+8 more)

### Community 52 - "model"
Cohesion: 0.22
Nodes (7): formatType(), listenForEvents(), syncViewport(), tickStats(), eventLine, eventMsg, model

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (38): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+30 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.10
Nodes (11): renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), assertViewFits(), TestViewsFitReportedTerminalSize(), RunFileEditorSession(), RunSession(), noticeModel (+3 more)

### Community 55 - "TrustGuard"
Cohesion: 0.06
Nodes (17): exeRow, tempGrant, tempMaskOp, tempRow, trustHook, trustRows, setSupersedeTrust(), liveTrust() (+9 more)

### Community 56 - "system.go"
Cohesion: 0.06
Nodes (48): applyEdit(), checkAndApply(), editAgain(), putConfig(), TestDaemonUnitWithMetadataOutput(), updateCatalogConfig(), dupFD(), installConfig() (+40 more)

### Community 57 - "absPath"
Cohesion: 0.12
Nodes (3): inNosuidNamespace(), IntegrationSuite, absPath()

### Community 58 - "engine"
Cohesion: 0.05
Nodes (26): engine, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), pathWithin(), TestPathWithinMatchesPrefixRule() (+18 more)

### Community 59 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.18
Nodes (10): findButton, span, column(), findAll(), Finder, lowerRunes(), newFindInput(), overlay() (+2 more)

### Community 60 - "ResolveTempBinary"
Cohesion: 0.16
Nodes (13): binaryStatOf(), ResolveTempBinary(), TestResolveTempBinaryClassesRuntime(), journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals() (+5 more)

### Community 61 - "fileedit_symlink_test.go"
Cohesion: 0.40
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.15
Nodes (19): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+11 more)

### Community 63 - "GuardEvent"
Cohesion: 0.10
Nodes (20): holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents() (+12 more)

### Community 64 - "expandPlaceholders"
Cohesion: 0.14
Nodes (8): symlinkTargetGlob(), expandPlaceholders(), TestExpandWhitelist(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 65 - "pinstate.go"
Cohesion: 0.10
Nodes (34): configureDaemonLogging(), confirmMasterKeyOverwrite(), resolveConfigPath(), runGenKey(), runLockdown(), runOneShotMode(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+26 more)

### Community 66 - "gui/session.go"
Cohesion: 0.05
Nodes (29): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+21 more)

### Community 68 - "rootHandle"
Cohesion: 0.20
Nodes (4): rootHandle, openRoot(), TestNoteRootMountFollowsMountID(), TestReadOnlyWalkToleratesDanglingSymlinks()

### Community 69 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 72 - "KernelDev"
Cohesion: 0.09
Nodes (27): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+19 more)

### Community 73 - "globBuilder"
Cohesion: 0.22
Nodes (8): entryWriters(), reserveChain(), underResource(), globBuilder, ParseGlobName(), TestGlobReservationsLibTrusted(), TestParseGlobName(), TestCatalogTrustGlobsAreReservable()

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "SetSectionWhitelist"
Cohesion: 0.07
Nodes (36): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, findSectionEnd(), findSectionStart(), isLibDirectiveText(), isPathDirectiveText() (+28 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 78 - "highlight.go"
Cohesion: 0.10
Nodes (7): diffModel, editorModel, Highlighter, lexerFor(), shebangLexer(), tokenClassOf(), changelogModel

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 81 - "hashCache"
Cohesion: 0.18
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 82 - "CandidateDir"
Cohesion: 0.21
Nodes (12): Discover(), DiscoverSystem(), CandidateDir, TestDiscoverOnlyExisting(), TestDiscoverSystemOnlyAbsolute(), TestPathsFor(), catalogEntryMatch(), findCatalogRoot() (+4 more)

### Community 83 - "runUninstall"
Cohesion: 0.19
Nodes (18): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), Delete(), Find(), find(), Backup (+10 more)

### Community 85 - "netModel"
Cohesion: 0.13
Nodes (16): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+8 more)

### Community 86 - "audit.go"
Cohesion: 0.16
Nodes (20): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+12 more)

### Community 87 - "TempBinary"
Cohesion: 0.14
Nodes (15): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+7 more)

### Community 88 - ".refreshOnce"
Cohesion: 0.16
Nodes (14): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), grantNewBinaries(), Lines(), TestLinesKeyOnLinkAndIncludePending() (+6 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "time.Duration"
Cohesion: 0.12
Nodes (10): drain(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity() (+2 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.17
Nodes (14): drainEphemeral(), admitEvent(), foldable(), targetProgram(), gateLogEntry, gateLogLimiter, listenForDaemonEvents(), NewDaemonModel() (+6 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.10
Nodes (10): foreignPinnedLink(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), parseProcs(), TestParseProcs(), TestReasonText() (+2 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 100 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "Pulse"
Cohesion: 0.28
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "NetGuardEvent"
Cohesion: 0.11
Nodes (17): NetGuardEvent, NetworkGuardRepository, NewMonitorUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo() (+9 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.09
Nodes (15): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule(), runHeadless(), runMonitor() (+7 more)

### Community 110 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.20
Nodes (3): mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 111 - "controlManager"
Cohesion: 0.25
Nodes (3): newControlManager(), startControlManager(), controlManager

### Community 112 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "eventUnitTest"
Cohesion: 0.27
Nodes (4): BpfEvent, eventUnitTest, Cstr(), DecodeBpfEvent()

### Community 116 - "hashFileShared"
Cohesion: 0.22
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 121 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "TestBatchInodeRows"
Cohesion: 0.50
Nodes (5): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), TestCountInodeRows()

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "writeEvent"
Cohesion: 0.33
Nodes (8): UIDResolver, NewUIDResolver(), writeEvent(), allowedEvent(), deniedEvent(), TestWriteEventAllowed(), TestWriteEventAllowedSuppressed(), TestWriteEventBlocked()

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "loadDaemonConfig"
Cohesion: 0.22
Nodes (9): loadDaemonConfig(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds(), SetUserHomes(), TestResolveConfinedRootOwnedParents() (+1 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 144 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 145 - "buildOneGuard"
Cohesion: 0.08
Nodes (37): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), newPinGeneration(), reloadOnce(), reloadSlotsNeeded() (+29 more)

### Community 147 - "render.go"
Cohesion: 0.12
Nodes (13): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 148 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 149 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 150 - "filevault.go"
Cohesion: 0.07
Nodes (46): harnessSuite, migrationTarget, classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord() (+38 more)

### Community 151 - "Monitor"
Cohesion: 0.15
Nodes (4): NewMonitor(), newMountTable(), Monitor, watchedInode

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 154 - "checkBtrfsKeys"
Cohesion: 0.33
Nodes (6): checkBtrfsKeys(), ensureBtrfsLayout(), FirstOnBtrfs(), OnBtrfs(), SuperblockDevPath(), TestSuperblockDevMatchesStatOffBtrfs()

### Community 155 - "io.Reader"
Cohesion: 0.25
Nodes (6): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype(), progressReader

### Community 156 - "find_test.go"
Cohesion: 0.48
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 157 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.06
Nodes (48): VettedInode, Multicall, MulticallKind, ExeKey(), IsMulticallRefusal(), multicallFamily(), MulticallRefusal(), RefuseMulticallLines() (+40 more)

### Community 161 - "RawTarget"
Cohesion: 0.24
Nodes (11): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+3 more)

### Community 163 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 164 - "runNetworkGuard"
Cohesion: 0.14
Nodes (20): CheckBPFLSM(), CheckEBPF(), ParseNetEventsFlag(), prepareDaemonStart(), runBPFCheck(), computeGuardBinaries(), confirmUnsafe(), discoverInfra() (+12 more)

### Community 165 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 166 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 168 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 179 - "EventType"
Cohesion: 0.09
Nodes (14): bpfGuardEvent, EventType, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel() (+6 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 184 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 185 - "SanitizeText"
Cohesion: 0.10
Nodes (18): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), systemPatterns(), systemRule(), runHeadless(), runHeadless(), binaryVetter (+10 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

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

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 453 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **32 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Process inspectors` connect `newWatchPattern` to `Daemon configuration`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13145444518619567 - nodes in this community are weakly interconnected._
- **Why does `Daemon configuration` connect `Daemon configuration` to `newWatchPattern`, `README.md`?**
  _High betweenness centrality (0.048) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.030303030303030304 - nodes in this community are weakly interconnected._
- **Why does `watchPattern` connect `newWatchPattern` to `time.Time`, `time.Duration`, `go_pkg_path_filepath`?**
  _High betweenness centrality (0.047) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.05202547080341417 - nodes in this community are weakly interconnected._