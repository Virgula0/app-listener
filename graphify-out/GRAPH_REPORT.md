# Graph Report - app-listener  (2026-10-17)

## Corpus Check
- 443 files · ~1,623,369 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 5664 nodes · 24023 edges · 199 communities (168 shown, 31 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2157 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5e7ff4ba`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- unit.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- mkdirs
- codeedit.go
- auditGroupCoverage
- installUI
- NetworkMonitor
- IntegrationSuite
- Network
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- runNetworkGuard
- guard.bpf.c
- offline.go
- vaultFS
- icon
- selection.go
- UI
- Serve
- IntegrationSuite
- github.com/charmbracelet/bubbles/textarea.Model
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- Guard
- sync.Mutex
- fileedit_model_test.go
- os.File
- IntegrationSuite
- absPath
- monitorUnitTest
- hashFileShared
- update/update.go
- SanitizeText
- stdio
- time.Duration
- runDaemon
- LibraryClosure
- GenerateConf
- fileedit_symlink_test.go
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- install/install.go
- update_test.go
- GuardInodeKey
- inode_dev
- buildOneGuard
- .Refresh
- exe_supersede.h
- check_watched
- runView
- bpfstats/main.go
- wm.go
- RawTarget
- integrationtests/networkguard_test.go
- FileTarget
- fcntl
- copyRegularAt
- BinaryEntry
- time.Time
- btrfsLayoutFrom
- changelog_test.go
- NewEventFanout
- RefreshCatalog
- NetEventType
- process_vm_readv.c
- GuardEvent
- fakeGuardRepo
- mayUpdate
- uninstall_system.go
- filevault_test.go
- trustManager
- Finder
- IntegrationSuite
- catalogrefresh.go
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- inode_dev.h
- inodeChain
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- CandidateDir
- policyEditor
- NewDumpHook
- Well-known bypass classes
- runMonitor
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- golang.org/x/sys/unix.Stat_t
- Build & sign release assets (reusable workflow)
- installPage
- KernelDev
- mountTable
- NetGuardEvent
- app-listener Terminal Demo Recording
- IntegrationSuite
- IntegrationSuite
- changelog.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- DaemonEvent
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- .FilterExistingWhitelist
- make test-integration (rootful Docker suite)
- fakeDaemon
- supersede_probe.c
- newTestVetter
- trace-app-libs.sh
- IntegrationSuite
- configedit_test.go
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- io.Reader
- runNetworkMonitor
- parseGuardEvent
- openBunTmp
- net_tester/main.go
- edit-protected
- Monitor
- fscrypt.go
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- EventType
- layout.go
- fetchReleases
- newTestTarget
- Installation
- .policyTools
- graphify-refresh.sh
- ParseSectionWhitelist
- Execute
- TestMain
- AddServeFlags
- feed.go
- shell
- launch_rule_at
- WithBottomBar
- writeWithin
- IntegrationSuite
- Limitations
- Troubleshooting
- fakeVault
- Vault
- Daemon configuration
- ParseEventsFlag
- Compatibility
- How it works
- tempGrantSetup
- bleedLayout
- dialLiveSession
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- .runEditorHarness
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- mountsUnitTest
- github.com/stretchr/testify/suite.Suite
- monitor: see what touches a path
- kernel-prof.sh
- pathCache
- .eventTypeAllowed
- .populateWatchInodes
- Desktop GUI

## God Nodes (most connected - your core abstractions)
1. `Guard` - 117 edges
2. `Config` - 96 edges
3. `GuardInodeKey` - 96 edges
4. `Load()` - 90 edges
5. `IntegrationSuite` - 82 edges
6. `absPath()` - 81 edges
7. `fileEditModel` - 77 edges
8. `shell` - 72 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `nextConfig()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/editprotected/editconfig.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (199 total, 31 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.11
Nodes (16): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, textAnswer, btrfsIoctlFsInfoArgs, multicallCacheKey (+8 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (96): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestDiffCatalogFlagRegistered(), TestIntegrationSuite(), TestMovePeer() (+88 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "unit.go"
Cohesion: 0.11
Nodes (23): needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), TestDaemonUnitWithMetadataOutput(), Mode, Request (+15 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.11
Nodes (34): code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image(), fill_inode_key(), get_task_exe_key() (+26 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "shellrc.go"
Cohesion: 0.07
Nodes (70): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), deploy() (+62 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (20): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+12 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (40): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), newWatchPattern(), systemBinaryPatterns(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped() (+32 more)

### Community 10 - "codeedit.go"
Cohesion: 0.11
Nodes (41): class, classOf(), column(), dedent(), Edit(), leadingIndent(), lines(), moveTo() (+33 more)

### Community 11 - "auditGroupCoverage"
Cohesion: 0.20
Nodes (12): auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+4 more)

### Community 12 - "installUI"
Cohesion: 0.15
Nodes (9): installUI, multiAnswer, ask(), diffText(), orDefault(), promptBody(), MultiSelectKeymap(), Choice (+1 more)

### Community 13 - "NetworkMonitor"
Cohesion: 0.13
Nodes (8): appletTag(), MulticallNameTags(), TestAppletTagStablePerName(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 15 - "Network"
Cohesion: 0.06
Nodes (21): binaryList, newBinaryList(), netGuardPill(), newActivity(), newExeResolver(), guardRow(), TestRingKeepsNewest(), clonePolicy() (+13 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (9): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), guardEventTypes(), le64HexKey(), parseGuardEvents(), guardEvent (+1 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (43): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+35 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (16): libBinaries(), planUpdaters(), resourceBinaries(), resourceLinks(), TestPlanUpdaters_LinkOwnsItsResourceBit(), Resource, GuardRepository, concurrencyLimit() (+8 more)

### Community 23 - "auth.go"
Cohesion: 0.06
Nodes (49): auditAfterEdit(), Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode() (+41 more)

### Community 24 - "runNetworkGuard"
Cohesion: 0.12
Nodes (21): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), WithSuspectsFrom(), DiscoverInfraBinaries() (+13 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (53): add_inode_to_guard(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with(), event_is_read_class(), evict_inode_from_guard() (+45 more)

### Community 26 - "offline.go"
Cohesion: 0.07
Nodes (39): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), TestCopyXattrsPreservesUserAttributes(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector() (+31 more)

### Community 27 - "vaultFS"
Cohesion: 0.14
Nodes (13): checkFreshDir(), openDirNoFollow(), renewDir(), TestRuntimeClassMarkers(), applyNewFileMeta(), defaultNewFileMode(), openRootT(), TestDefaultNewFileMode() (+5 more)

### Community 28 - "icon"
Cohesion: 0.07
Nodes (35): binaryRow, checkGrid, filesPage, homePage, iconTile, panel, pill, statusTile (+27 more)

### Community 29 - "selection.go"
Cohesion: 0.11
Nodes (18): pos, Highlighter, lexerFor(), shebangLexer(), cellToRune(), clipboardText(), cursorPos(), Selection (+10 more)

### Community 30 - "UI"
Cohesion: 0.10
Nodes (19): showChangelog(), confirmUpdate(), isTerminal(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun(), cancelUI (+11 more)

### Community 31 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 32 - "IntegrationSuite"
Cohesion: 0.26
Nodes (3): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite

### Community 33 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.16
Nodes (13): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.24
Nodes (12): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials(), TestSetupDumpLogRefusesSymlink() (+4 more)

### Community 37 - "Guard"
Cohesion: 0.10
Nodes (6): rootHandle, SharedPinDegraded(), Guard, patternMatches(), releaseRootTrust(), TestPatternMatchesExpandOnlyRootPlacedDirs()

### Community 38 - "sync.Mutex"
Cohesion: 0.10
Nodes (20): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines(), startControlServer() (+12 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.11
Nodes (37): OpenRegularNoFollow(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed() (+29 more)

### Community 40 - "os.File"
Cohesion: 0.05
Nodes (56): setConfinedHomes(), VettedInode, Multicall, MulticallKind, WithVettedKeys(), ExeKey(), IsMulticallRefusal(), multicallFamily() (+48 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+7 more)

### Community 42 - "absPath"
Cohesion: 0.08
Nodes (4): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, absPath()

### Community 44 - "hashFileShared"
Cohesion: 0.19
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 45 - "update/update.go"
Cohesion: 0.11
Nodes (28): applyUpdate(), assetsFor(), channelOf(), checkUpdatePreconditions(), confirmUnit(), deploy(), downloadAndVerify(), downloadReleaseFiles() (+20 more)

### Community 46 - "SanitizeText"
Cohesion: 0.05
Nodes (28): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), newHashCache() (+20 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "time.Duration"
Cohesion: 0.14
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 49 - "runDaemon"
Cohesion: 0.07
Nodes (23): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), newControlManager(), startControlManager() (+15 more)

### Community 50 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 51 - "GenerateConf"
Cohesion: 0.05
Nodes (64): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), askEncryption(), editConfig(), libraryBlocksFromCandidates(), runConfigEditor(), sectionsFromCandidates() (+56 more)

### Community 52 - "fileedit_symlink_test.go"
Cohesion: 0.35
Nodes (10): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+2 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.11
Nodes (47): systemPatterns(), systemRule(), grantNewBinaries(), watchGroup, Line(), Lines(), TestLinesKeyOnLinkAndIncludePending(), addResource() (+39 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (47): diffModel, editorModel, Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace() (+39 more)

### Community 55 - "TrustGuard"
Cohesion: 0.06
Nodes (18): tempGrant, trustHook, trustRows, deleteResKeys(), setSupersedeTrust(), liveTrust(), loadPinnedExeMaps(), stripPinnedAllow() (+10 more)

### Community 56 - "install/install.go"
Cohesion: 0.06
Nodes (62): reloadDaemonForPasswordChange(), promptEditPassword(), applyLiveRefresh(), buildBinaryIfNeeded(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), encryptDirectories() (+54 more)

### Community 57 - "update_test.go"
Cohesion: 0.12
Nodes (11): rsaPublicKeyPEM(), signedChecksum(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled(), TestNewerThanStable(), TestParseChecksum(), TestParsePublicKey() (+3 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (53): engine, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), TestSetInspectorsBounded(), appendInodeRow() (+45 more)

### Community 59 - "inode_dev"
Cohesion: 0.18
Nodes (15): exe_is_superseded(), exe_refused(), inode_dev(), check_and_emit_args(), code_suspect(), current_is_inspector(), exe_tag_bits(), fill_path() (+7 more)

### Community 60 - "buildOneGuard"
Cohesion: 0.06
Nodes (52): backingDeviceUnion(), buildConcurrency(), buildGuards(), buildOneGuard(), eventFilterOptions(), reloadOnce(), reloadSlotsNeeded(), startGuardedDaemon() (+44 more)

### Community 61 - ".Refresh"
Cohesion: 0.04
Nodes (13): appTheme, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, asset() (+5 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 63 - "check_watched"
Cohesion: 0.24
Nodes (14): admit_refusal(), check_watched(), emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg() (+6 more)

### Community 64 - "runView"
Cohesion: 0.08
Nodes (12): logHook, logLine, runHost, runView, stageItem, stageState, logImportance(), captureLogs() (+4 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.05
Nodes (62): CheckBPFLSM(), CheckEBPF(), configureDaemonLogging(), newPinGeneration(), prepareDaemonStart(), resolveConfigPath(), runBPFCheck(), runLockdown() (+54 more)

### Community 66 - "wm.go"
Cohesion: 0.06
Nodes (44): invokingIDs(), TestWriteUserFileRefusesSymlinks(), writeUserFile(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv() (+36 more)

### Community 67 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 68 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 69 - "FileTarget"
Cohesion: 0.12
Nodes (13): Target, specFromForm(), guardPill(), specSummary(), EventTypes(), Files, FileSpec, FileTarget (+5 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "copyRegularAt"
Cohesion: 0.09
Nodes (23): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), copyDirContents(), copyEntry(), copyRegularAt() (+15 more)

### Community 72 - "BinaryEntry"
Cohesion: 0.06
Nodes (24): deferredBinary, exeRow, tempMaskOp, tempRow, snapNet, BinaryEntry, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), checkBtrfsKeys() (+16 more)

### Community 73 - "time.Time"
Cohesion: 0.16
Nodes (9): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, sweepID, sweepMemo (+1 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 76 - "NewEventFanout"
Cohesion: 0.24
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 77 - "RefreshCatalog"
Cohesion: 0.13
Nodes (18): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+10 more)

### Community 78 - "NetEventType"
Cohesion: 0.08
Nodes (15): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), Cstr(), FormatAddr() (+7 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "GuardEvent"
Cohesion: 0.07
Nodes (16): holdHeadless(), HeadlessLine(), runGuardHeadless(), snapGuard, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents() (+8 more)

### Community 82 - "mayUpdate"
Cohesion: 0.20
Nodes (11): hardLinked(), mayUpdate(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary(), TestPlanUpdaters_ScopedToOwnResources() (+3 more)

### Community 83 - "uninstall_system.go"
Cohesion: 0.05
Nodes (63): collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), runMaintenanceMode() (+55 more)

### Community 84 - "filevault_test.go"
Cohesion: 0.11
Nodes (41): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+33 more)

### Community 85 - "trustManager"
Cohesion: 0.10
Nodes (24): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+16 more)

### Community 86 - "Finder"
Cohesion: 0.14
Nodes (19): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+11 more)

### Community 88 - "catalogrefresh.go"
Cohesion: 0.33
Nodes (3): isAbsoluteCandidate(), TestCatalogSanity(), hint

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "inode_dev.h"
Cohesion: 0.24
Nodes (8): ctx_ptr(), sb_dev(), btrfs_copy_gate(), guard_file_ioctl(), guard_file_ioctl_compat(), guard_inode_free(), is_btrfs_copy_ioctl(), trust_inode_free()

### Community 92 - "inodeChain"
Cohesion: 0.13
Nodes (23): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), openRoot(), physicalParent(), statFD() (+15 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.09
Nodes (12): foreignPinnedLink(), reasonText(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs() (+4 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "CandidateDir"
Cohesion: 0.05
Nodes (37): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, GlobChild, statFunc (+29 more)

### Community 100 - "policyEditor"
Cohesion: 0.12
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 101 - "NewDumpHook"
Cohesion: 0.24
Nodes (11): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "runMonitor"
Cohesion: 0.28
Nodes (6): MakeDisplayPaths(), runHeadless(), runMonitor(), runTUI(), NewModel(), MonitorUseCase

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.06
Nodes (26): binSource, installAction, installedBinary, InstallMode, installPage, InstallRequest, latestRelease, Options (+18 more)

### Community 110 - "KernelDev"
Cohesion: 0.08
Nodes (32): BinaryStat, binaryStatOf(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs() (+24 more)

### Community 111 - "mountTable"
Cohesion: 0.22
Nodes (8): MountInfo, TestNoteRootMountFollowsMountID(), MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), mountTable

### Community 112 - "NetGuardEvent"
Cohesion: 0.09
Nodes (19): snapNetGuard, NetGuardEvent, NetworkGuardRepository, NewGuardUseCase(), NewMonitorUseCase(), NetworkGuardUseCase, NewNetworkGuardUseCase(), NewNetworkMonitorUseCase() (+11 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 116 - "changelog.go"
Cohesion: 0.27
Nodes (16): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+8 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "DaemonEvent"
Cohesion: 0.08
Nodes (37): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+29 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - ".FilterExistingWhitelist"
Cohesion: 0.11
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.13
Nodes (15): forwardRequest, TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess (+7 more)

### Community 143 - "io.Reader"
Cohesion: 0.20
Nodes (7): downloadFile(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), BtrfsMounted(), mountinfoHasFstype(), TestMountinfoHasFstype(), progressReader

### Community 144 - "runNetworkMonitor"
Cohesion: 0.18
Nodes (8): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NewNetModel(), NetworkMonitorUseCase

### Community 145 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 146 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 147 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 148 - "edit-protected"
Cohesion: 0.29
Nodes (7): edit-protected, Editing the daemon config, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 149 - "Monitor"
Cohesion: 0.14
Nodes (3): NewMonitor(), Monitor, watchedInode

### Community 150 - "fscrypt.go"
Cohesion: 0.06
Nodes (44): confirmMasterKeyOverwrite(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runGenKey(), deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget() (+36 more)

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.12
Nodes (34): loadDaemonConfig(), markReady(), noteReloadedSettings(), notifySystemdReady(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), inodeBatchFn, fakeBatches() (+26 more)

### Community 154 - "EventType"
Cohesion: 0.08
Nodes (12): snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+4 more)

### Community 155 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 156 - "fetchReleases"
Cohesion: 0.29
Nodes (7): fetchRelease(), fetchReleases(), FetchVerified(), TestFetchReleaseByTag(), TestFetchReleases(), TestFetchReleasesHTTPError(), updateHTTPClient()

### Community 157 - "newTestTarget"
Cohesion: 0.11
Nodes (20): NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst() (+12 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 161 - "ParseSectionWhitelist"
Cohesion: 0.29
Nodes (7): isLibDirectiveText(), ParseSectionHeaderPath(), unquotePath(), IsLibraryDirective(), parseLibDirLine(), ParseSectionWhitelist(), TestParseSectionWhitelist()

### Community 162 - "Execute"
Cohesion: 0.20
Nodes (10): installHooks(), releaseBuild(), Requested(), Run(), TestRequested(), TestSubcommandOutranksGUI(), Execute(), initLogger() (+2 more)

### Community 163 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 164 - "AddServeFlags"
Cohesion: 0.33
Nodes (6): AddServeFlags(), init(), init(), init(), init(), init()

### Community 165 - "feed.go"
Cohesion: 0.07
Nodes (37): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+29 more)

### Community 166 - "shell"
Cohesion: 0.07
Nodes (24): pageDef, pageID, resizeWatch, shell, sizedDialog, newLogo(), newShell(), Run() (+16 more)

### Community 167 - "launch_rule_at"
Cohesion: 0.50
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 168 - "WithBottomBar"
Cohesion: 0.24
Nodes (7): decryptDirectories(), TestDecryptDirectoriesEmpty(), TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 170 - "writeWithin"
Cohesion: 0.22
Nodes (9): atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+1 more)

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 176 - "Daemon configuration"
Cohesion: 0.20
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 178 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "tempGrantSetup"
Cohesion: 0.20
Nodes (11): tempJournalRows(), journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure() (+3 more)

### Community 185 - "dialLiveSession"
Cohesion: 0.22
Nodes (11): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+3 more)

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

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

### Community 195 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.18
Nodes (4): harnessSuite, mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 200 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 206 - "Desktop GUI"
Cohesion: 0.18
Nodes (11): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Tables (+3 more)

## Knowledge Gaps
- **32 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+27 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 523 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **31 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Editor keys` connect `auth.go` to `uninstall_system.go`, `edit-protected`?**
  _High betweenness centrality (0.052) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `open_by_handle_at POC` to the rest of the system?**
  _32 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11012345679012346 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`, `auth.go`?**
  _High betweenness centrality (0.052) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.027829313543599257 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `daemon/daemon.go`, `IntegrationSuite`, `TestMain`, `integrationtests/networkguard_test.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `auth.go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._